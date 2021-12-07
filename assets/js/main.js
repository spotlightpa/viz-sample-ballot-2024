import Alpine from "alpinejs/src/index.js";

import { initFrameAndPoll } from "@newswire/frames";

import { onLoad } from "./utils/dom-utils.js";
import { addGAListeners } from "./utils/google-analytics.js";

import * as L from "leaflet";

L.Marker.prototype.options.icon = L.icon({
  iconUrl: "/images/marker-icon.png",
  iconRetinaUrl: "/images/marker-icon-2x.png",
  shadowUrl: "/images/marker-shadow.png",
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  tooltipAnchor: [16, -28],
  shadowSize: [41, 41],
});

import * as params from "@params";

const locate = () =>
  new Promise((resolve, reject) =>
    navigator.geolocation.getCurrentPosition(resolve, reject)
  );

Alpine.directive(
  "template",
  (el, { expression }, { effect, evaluateLater }) => {
    let evalStr = expression
      ? "`" + expression + "`"
      : "`" + el.innerHTML.trim() + "`";
    let evaluate = evaluateLater(evalStr);

    effect(() => {
      evaluate((value) => {
        el.innerHTML = value;
      });
    });
  }
);

let commaFormatter = new Intl.NumberFormat("en-US");
Alpine.magic("comma", () => (n) => commaFormatter.format(n));

Alpine.store("state", {
  oldHouse: "103",
  newHouse: "103",
  oldSenate: "15",
  newSenate: "15",

  lat: 40.26444,
  long: -76.88375,

  get latLong() {
    return [this.lat, this.long];
  },

  async updateAddress() {
    if (!this.address) {
      return;
    }
    return this.callAPI(
      "/api/by-address?address=" + encodeURIComponent(this.address)
    );
  },

  async updateLocation([lat, long]) {
    this.lat = lat;
    this.long = long;

    return this.callAPI(
      "/api/by-location?lat=" +
        encodeURIComponent(this.lat) +
        "&long=" +
        encodeURIComponent(this.long)
    );
  },

  async callAPI(path) {
    return fetch(params.apiBaseURL + path)
      .then((rsp) => rsp.json())
      .then((json) => {
        if (json.error_message) {
          throw Error(json.error_message);
        }
        this.update(json);
      });
  },

  update(json) {
    let { address, lat, long, old_house, new_house, old_senate, new_senate } =
      json;

    this.address = address || this.address;
    this.lat = lat || this.lat;
    this.long = long || this.long;
    this.oldHouse = old_house || this.oldHouse;
    this.newHouse = new_house || this.newHouse;
    this.oldSenate = old_senate || this.oldSenate;
    this.newSenate = new_senate || this.newSenate;
    if (
      ![old_house, new_house, old_senate, new_senate].every((item) => !!item)
    ) {
      throw new Error("District could not be found.");
    }
  },
});

Alpine.data("app", () => {
  return {
    error: null,

    isLoading: false,
    recentlyChanged: false,

    init() {
      let timeoutID = null;
      this.$watch("isLoading", (val) => {
        if (val) {
          this.recentlyChanged = true;
          window.clearTimeout(timeoutID);
          timeoutID = window.setTimeout(() => {
            this.recentlyChanged = false;
          }, 1000);
        }
      });
    },

    get isLoadingThrottled() {
      return this.isLoading || this.recentlyChanged;
    },

    async byAddress() {
      if (this.isLoading) {
        return;
      }
      this.isLoading = true;
      return this.$store.state
        .updateAddress()
        .then(() => {
          this.error = null;
        })
        .catch((e) => {
          this.error = e;
        })
        .finally(() => {
          this.isLoading = false;
        });
    },

    async byLocation() {
      if (this.isLoading) {
        return;
      }
      try {
        let { coords } = await locate();
        this.isLoading = true;
        await this.$store.state.updateLocation([
          coords.latitude,
          coords.longitude,
        ]);
        this.error = null;
      } catch (e) {
        this.error = e;
      } finally {
        this.isLoading = false;
      }
    },
  };
});

Alpine.data("map", ({ age, kind }) => {
  return {
    age,
    kind,

    map: null,
    latLong: null,
    marker: null,
    layer: null,
    geojson: "",
    props: null,

    init() {
      this.map = L.map(this.$refs.leaflet, {
        center: [this.$store.state.lat, this.$store.state.long],
        scrollWheelZoom: false,
        zoom: 12,
      });
      this.map.createPane("labels");
      this.map.getPane("labels").style.zIndex = 400;
      this.map.getPane("labels").style.pointerEvents = "none";
      L.tileLayer(
        "https://{s}.basemaps.cartocdn.com/rastertiles/voyager_nolabels/{z}/{x}/{y}{r}.png",
        {
          attribution:
            '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
          subdomains: "abcd",
          maxZoom: 12,
        }
      ).addTo(this.map);
      L.tileLayer(
        "https://{s}.basemaps.cartocdn.com/rastertiles/voyager_only_labels/{z}/{x}/{y}{r}.png",
        {
          attribution:
            '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
          subdomains: "abcd",
          maxZoom: 12,
          pane: "labels",
        }
      ).addTo(this.map);

      const watchDrag = () => {
        let { lat, lng } = this.marker.getLatLng();
        if (lat && lng) {
          this.layer.remove();
          return this.$store.state.updateLocation([lat, lng]);
        }
      };

      this.$watch("latLong", (latLong) => {
        if (this.marker) this.marker.remove();

        this.marker = L.marker(latLong, { draggable: true, zIndexOffset: 651 });
        this.marker.on("moveend", watchDrag);
        this.marker.addTo(this.map);

        this.map.flyTo(latLong);
      });

      this.$watch("geojson", async (url) => {
        let geojsonFeature;
        try {
          geojsonFeature = await fetch(url).then((rsp) => rsp.json());
        } catch (e) {
          this.error = e;
          return;
        }
        if (this.layer) this.layer.remove();
        let layer = L.geoJSON(geojsonFeature);
        layer.setStyle({
          weight: 1,
          color: "#cda52d",
          fillColor: "#ffcb05",
          fillOpacity: 0.5,
        });
        this.layer = layer;
        this.props = geojsonFeature.features[0].properties;
        layer.bindPopup(() => `District ${this.props.district}`);
        layer.addTo(this.map);
        this.map.flyToBounds(layer.getBounds());
      });
    },
  };
});

Alpine.start();

initFrameAndPoll();

onLoad(() => {
  addGAListeners();
});

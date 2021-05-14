import Alpine from "alpinejs/src/index.js";

import { initFrameAndPoll } from "@newswire/frames";

import { onLoad } from "./utils/dom-utils.js";
import { addGAListeners } from "./utils/google-analytics.js";

import * as L from "leaflet";

import * as params from "@params";

initFrameAndPoll();

onLoad(() => {
  addGAListeners();
});

const locate = () =>
  new Promise((resolve, reject) =>
    navigator.geolocation.getCurrentPosition(resolve, reject)
  );

Alpine.data("map", () => {
  return {
    lat: null,
    long: null,
    isLoading: false,
    recentlyChanged: false,
    error: null,
    address: "",

    init() {
      this.$watch("isLoading", (val) => {
        if (val) {
          this.recentlyChanged = true;
          window.setTimeout(() => {
            this.recentlyChanged = false;
          }, 1000);
        }
      });
    },

    get isLoadingThrottled() {
      return this.isLoading || this.recentlyChanged;
    },

    async geolocate() {
      if (this.isLoading || !this.address) {
        return;
      }
      this.isLoading = true;
      let { address } = this;
      if (!/PA|Pennsylvania/.exec(address)) {
        address += ", PA";
      }
      await fetch(
        `https://maps.googleapis.com/maps/api/geocode/json?address=${address}&key=${params.apiKey}`
      )
        .then((rsp) => rsp.json())
        .then((json) => {
          if (json.error_message) {
            throw Error(json.error_message);
          }
          this.address = json.results[0].formatted_address;
          let { lat, lng } = json.results[0].geometry.location;
          this.lat = lat;
          this.long = lng;
        })
        .catch((e) => {
          this.error = e;
        })
        .finally(() => {
          this.isLoading = false;
        });
    },

    async locate() {
      try {
        let { coords } = await locate();
        this.lat = coords.latitude;
        this.long = coords.longitude;
      } catch (e) {
        this.error = e;
      }
    },

    map() {
      L.map(this.$el);
      // map.createPane("labels");

      // var CartoDB_VoyagerLabelsUnder = L.tileLayer(
      //   "https://{s}.basemaps.cartocdn.com/rastertiles/voyager_labels_under/{z}/{x}/{y}{r}.png",
      //   {
      //     attribution:
      //       '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
      //     subdomains: "abcd",
      //     maxZoom: 20,
      //   }
      // );
    },
  };
});

Alpine.start();

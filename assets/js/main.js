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

Alpine.data("app", () => {
  return {
    lat: 40.26444,
    long: -76.88375,
    oldHouse: "103",
    newHouse: "103",
    oldSenate: "15",
    newSenate: "15",
    address: "",
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
      if (this.isLoading || !this.address) {
        return;
      }
      return this.callAPI(
        "/api/by-address?address=" + encodeURIComponent(this.address)
      );
    },

    async byLocation() {
      if (this.isLoading) {
        return;
      }

      try {
        let { coords } = await locate();
        this.lat = coords.latitude;
        this.long = coords.longitude;
      } catch (e) {
        this.error = e;
        return;
      }
      return this.callAPI(
        "/api/by-location?lat=" +
          encodeURIComponent(this.lat) +
          "&long=" +
          encodeURIComponent(this.long)
      );
    },
    async callAPI(path) {
      this.isLoading = true;
      await fetch(params.apiBaseURL + path)
        .then((rsp) => rsp.json())
        .then((json) => {
          if (json.error_message) {
            throw Error(json.error_message);
          }
          this.error = null;
          let {
            address,
            lat,
            long,
            old_house,
            new_house,
            old_senate,
            new_senate,
          } = json;
          this.address = address || this.address;
          this.lat = lat || this.lat;
          this.long = long || this.long;
          this.oldHouse = old_house || this.oldHouse;
          this.newHouse = new_house || this.newHouse;
          this.oldSenate = old_senate || this.oldSenate;
          this.newSenate = new_senate || this.newSenate;
          if (
            ![old_house, new_house, old_senate, new_senate].every(
              (item) => !!item
            )
          ) {
            throw new Error("District could not be found.");
          }
        })
        .catch((e) => {
          this.error = e;
        })
        .finally(() => {
          this.isLoading = false;
        });
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

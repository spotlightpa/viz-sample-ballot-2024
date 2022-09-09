import Alpine from "alpinejs/src/index.js";
import { initFrameAndPoll } from "js/framer/index.mjs";

import { addGAListeners, reportClick } from "./utils/google-analytics.js";

import * as params from "@params";

Alpine.magic("report", () => (ev) => reportClick(ev));

const locate = () =>
  new Promise((resolve, reject) =>
    navigator.geolocation.getCurrentPosition(resolve, reject)
  );

async function callAPI(path, query = {}) {
  let url = `${params.apiBaseURL}${path}?${new URLSearchParams(query)}`;

  return fetch(url)
    .then((rsp) => rsp.json())
    .then((json) => {
      if (json.error_message) {
        throw Error(json.error_message);
      }
      return json;
    });
}

Alpine.data("app", () => {
  return {
    address: "",
    rawData: null,
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

    lookup(kind) {
      return this.rawData[kind] || [];
    },

    get isLoadingThrottled() {
      return this.isLoading || this.recentlyChanged;
    },

    async byAddress() {
      if (this.isLoading) {
        return;
      }
      this.isLoading = true;

      return callAPI("/api/candidates-by-address", {
        address: this.address,
      })
        .then((data) => {
          this.error = null;
          this.rawData = data;
          this.address = data.address;
        })
        .catch((e) => {
          this.error = e;
        })
        .finally(() => {
          this.isLoading = false;
          this.$report({ target: this.$el });
        });
    },

    async byLocation() {
      if (this.isLoading) {
        return;
      }
      try {
        let { coords } = await locate();
        this.isLoading = true;
        await callAPI("/api/candidates-by-location", {
          lat: coords.latitude,
          long: coords.longitude,
        }).then((data) => {
          this.error = null;
          this.rawData = data;
        });
      } catch (e) {
        this.error = e;
      } finally {
        this.isLoading = false;
        this.$report({ target: this.$el });
      }
    },
  };
});

Alpine.start();
addGAListeners();
initFrameAndPoll();

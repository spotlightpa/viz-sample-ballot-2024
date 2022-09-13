module.exports = {
  purge: false,
  theme: {
    fontFamily: {
      serif: "Merriweather, serif",
      sans:
        'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", ' +
        'Roboto, "Helvetica Neue", Arial, "Noto Sans", "Apple Color Emoji", ' +
        '"Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji", sans-serif',
      rale: "Raleway, Helvetica, Arial, sans-serif",
    },
    colors: {
      transparent: "transparent",
      current: "currentColor",
      white: "#fff",
      black: "#000",
      "g-1": "#fcfcfc",
      "g-2": "#fafafa",
      "g-3": "#f5f5f5",
      "g-4": "#efefef",
      "g-5": "#dbdbdb",
      "g-6": "#8a8a8a",
      "g-7": "#5d5d5d",
      "g-8": "#2d3748",
      "g-9": "#1a202c",
      beige: "#f4f1ee",
      robin: "#99d9f1",
      blue: "#009edb",
      "blue-darker": "#1982BD",
      "tw-blue": "#99d9f1",
      "fb-blue": "#009edb",
      orange: "#ff6c36",
      green: "#78bc20",
      yellow: "#ffcb05",
      "yellow-darker": "#cda52d",
      darkblue: "#22416e",
      goldenrod: "#fff1bd",
      cyan: "#e5f6ff",
      "chart-red": "#d7191c",
      "chart-orange": "#F59E0B",
      "chart-green": "#059669",
      "chart-blue": "#2b83ba",
      "chart-purple": "#7b3294",
      red: {
        0: "#fff0f0",
        1: "#ffdddd",
        2: "#ffc0c0",
        3: "#ff9494",
        4: "#ff5757",
        5: "#ff2323",
        DEFAULT: "#ff0000",
        6: "#ff0000",
        7: "#d70000",
        8: "#b10303",
        9: "#920a0a",
      },
      navy: {
        0: "#f3f6fc",
        1: "#e6edf8",
        2: "#c8d9ef",
        3: "#97b8e2",
        4: "#6093d0",
        5: "#3b75bc",
        6: "#2b5b9e",
        7: "#244980",
        DEFAULT: "#22416e",
        8: "#22416e",
        9: "#20375a",
      },
    },
    extend: {
      boxShadow: {
        beige: "0 0 0 3px #f4f1ee80",
      },
      dropShadow: {
        outline: "0 1px 1px rgba(0, 0, 0, 0.5)",
      },
      lineHeight: {
        normal: "1.6",
        hed: "1.15",
      },
      margin: {
        21: "5.25rem", // for staggered sidebars
      },
      spacing: {
        "16x9": `${(100 * 9) / 16}%`,
        "6x4": `${(100 * 4) / 6}%`,
      },
      screens: {
        sm: "640px",
      },
      maxWidth: {
        content: "730px",
      },
    },
  },
  variants: {
    extend: {
      boxShadow: ["active", "group-focus"],
      ringWidth: ["hover", "active"],
    },
  },
  plugins: [],
  future: {
    purgeLayersByDefault: true,
    removeDeprecatedGapUtilities: true,
  },
  darkMode: "media",
};

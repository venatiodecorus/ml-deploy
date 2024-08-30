import daisyui from "daisyui";
import typography from "@tailwindcss/typography";

module.exports = {
  content: ["./src/frontend/templates/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [daisyui, typography],
};

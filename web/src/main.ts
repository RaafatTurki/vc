import { mount } from "svelte"
import App from "./App.svelte"
import "./app.css"

mount(App, { target: document.querySelector("#app")! })

if ("serviceWorker" in navigator) {
  window.addEventListener("load", () => navigator.serviceWorker.register("/sw.js"))
}

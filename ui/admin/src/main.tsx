import { i18n } from "@allape/gocrud-react";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.scss";
import App from "./App.tsx";
import TranslationEn from "./i18n/en.ts";
import TranslationZh from "./i18n/zh.ts";

i18n
  .setup({
    zh: {
      translation: TranslationZh,
    },
    en: {
      translation: TranslationEn,
    },
  })
  .then(() => {
    createRoot(document.getElementById("root")!).render(
      <StrictMode>
        <App />
      </StrictMode>,
    );
  });

import { i18n } from "@allape/gocrud-react";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.scss";
import App from "./App.tsx";

i18n
  .setup({
    zh: {
      translation: {
        // ...i18n.ZHCN,
      },
    },
    en: {
      translation: {
        ...i18n.EN,

        id: "ID",
        unknown: "Unknown",
        select: "Select",
        createdAt: "Created At",
        updatedAt: "Updated At",

        tag: {
          _: "Tag",
          name: "Name",
          alias: "Alias",
          priority: "Priority",
          description: "Description",

          aliasDesc: "Separated by comma(,)",
          priorityDesc: "Larger for higher priority",
        },

        gallery: {
          _: "Gallery",
          isPublic: "Is Public",
          name: "Name",
          createdBy: "Created By",
          priority: "Priority",
          description: "Description",

          isPublicYesOrNo: {
            yes: "Public",
            no: "Private",
          },
        },
      },
    },
  })
  .then(() => {
    createRoot(document.getElementById("root")!).render(
      <StrictMode>
        <App />
      </StrictMode>,
    );
  });

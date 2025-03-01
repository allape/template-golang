import { i18n } from "@allape/gocrud-react";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.scss";
import App from "./App.tsx";

i18n
  .setup({
    zh: {
      translation: {
        ...i18n.ZHCN,

        id: "ID",
        unknown: "未知",
        select: "选择",
        createdAt: "创建时间",
        updatedAt: "更新时间",

        user: {
          _: "用户",
          name: "姓名",
          description: "描述",
        },
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

        user: {
          _: "User",
          name: "Name",
          description: "Description",
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

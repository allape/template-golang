import { i18n } from "@allape/gocrud-react";
import en from "./en.ts";

const Translation: typeof en = {
  ...i18n.ZHCN,

  id: "ID",
  unknown: "Unknown",
  select: "Select",
  createdAt: "Created At",
  updatedAt: "Updated At",

  home: {
    search: "搜索",
  },
  gallery: {
    galleryNotValid: "图册数据有误",
  },
};

export default Translation;

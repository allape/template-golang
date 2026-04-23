import { i18n } from "@allape/gocrud-react";

const Translation = {
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
    color: "Color",
    description: "Description",

    aliasDesc: "Separated by comma(,)",
    priorityDesc: "Larger for higher priority",
  },

  item: {
    _: "Item",
    name: "Name",
    src: "Source",
    priority: "Priority",
    description: "Description",

    galleries: "Galleries",
    tags: "Tags",

    continuesUpload: "Continues Upload",
    priorityDesc: "Larger for higher priority",
    srcRequired: "Please upload an image",
    saved: "Saved 🎉",
  },

  gallery: {
    _: "Gallery",
    isPublic: "Is Public",
    name: "Name",
    createdBy: "Created By",
    priority: "Priority",
    description: "Description",

    photos: "Photos",

    isPublicYesOrNo: {
      yes: "Public",
      no: "Private",
    },
  },
};

export default Translation;

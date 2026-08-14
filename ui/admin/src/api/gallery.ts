import Crudy, { AntdM2MConnectorHandler, config } from "@allape/gocrud-react";
import {
  IGallery,
  IGalleryItem,
  IGalleryTag,
} from "common/src/model/gallery.ts";
import { IItem } from "common/src/model/item.ts";
import { ITag } from "common/src/model/tag.ts";
import { IGallerySearchParams } from "../model/gallery.ts";
import { ItemCrudy } from "./item.ts";
import { TagCrudy } from "./tag.ts";

export const GalleryCrudy = new Crudy<IGallery, IGallerySearchParams>(
  `${config.SERVER_URL}/gallery`,
);

export const GalleryItemHandler = new AntdM2MConnectorHandler<
  IGallery,
  IItem,
  IGalleryItem
>(
  `${config.SERVER_URL}/gallery-item`,
  GalleryCrudy,
  ItemCrudy,
  "galleryId",
  "itemId",
);

export const GalleryTagHandler = new AntdM2MConnectorHandler<
  IGallery,
  ITag,
  IGalleryTag
>(
  `${config.SERVER_URL}/gallery-tag`,
  GalleryCrudy,
  TagCrudy,
  "galleryId",
  "tagId",
);

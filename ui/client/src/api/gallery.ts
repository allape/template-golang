import { Config } from "@allape/gocrud";
import { antdget, config } from "@allape/gocrud-react";
import { IGallery } from "common/src/model/gallery.ts";
import { IItem, IItemTag } from "common/src/model/item.ts";
import { ITag } from "common/src/model/tag.ts";

export interface IGalleryInfo {
  info: IGallery;
  tags?: ITag[];
  /**
   * @deprecated
   */
  cover?: IItem;
}

export function getAllGalleries(
  cfg?: Config<IGalleryInfo[]>,
): Promise<IGalleryInfo[]> {
  return antdget(`${config.SERVER_URL}/gallery/all`, cfg);
}

export interface IGalleryDetail {
  gallery: IGallery;
  items: IItem[];
  itemTags: IItemTag[];
  tags: ITag[];
}

export function getDetailById(
  id: IGallery["id"],
  cfg?: Config<IGalleryDetail>,
): Promise<IGalleryDetail> {
  return antdget(`${config.SERVER_URL}/gallery/detail/${id}`, cfg);
}

/**
 * @param galleryId
 * @param itemId retrieve the first image of gallery when 0
 * @param imageType
 */
export function toImageURL(
  galleryId: IGallery["id"],
  itemId: IItem["id"] = 0,
  imageType: "src" | "thumbnail" = "src",
): string {
  return `${config.SERVER_URL}/gallery/${imageType}/${galleryId}/${itemId}`;
}

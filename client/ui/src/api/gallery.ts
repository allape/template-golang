import { Config } from "@allape/gocrud";
import { antdget, config } from "@allape/gocrud-react";
import { IGallery } from "../model/gallery.ts";
import { IItem, IItemTag } from "../model/item.ts";
import { ITag } from "../model/tag.ts";

export interface IGalleryInfo {
  info: IGallery;
  tags: ITag[];
  /**
   * @deprecated
   */
  cover: IItem;
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
 */
export function toImageURL(
  galleryId: IGallery["id"],
  itemId: IItem["id"] = 0,
): string {
  // return `${config.STATIC_SERVER_URL}${item.src}`;
  return `${config.SERVER_URL}/gallery/image/${galleryId}/${itemId}`;
}

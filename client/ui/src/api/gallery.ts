import { config, get } from "@allape/gocrud-react";
import { IGallery } from "../model/gallery.ts";
import { IItem, IItemTag } from "../model/item.ts";
import { ITag } from "../model/tag.ts";

export function getAllGalleries(): Promise<IGallery[]> {
  return get(`${config.SERVER_URL}/gallery/all`);
}

export interface IGalleryDetailPayload {
  gallery: IGallery;
  items: IItem[];
  itemTags: IItemTag[];
  tags: ITag[];
}

export function getDetailById(
  id: IGallery["id"],
): Promise<IGalleryDetailPayload> {
  return get(`${config.SERVER_URL}/gallery/detail/${id}`);
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

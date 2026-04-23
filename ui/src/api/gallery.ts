import Crudy, { config, get } from "@allape/gocrud-react";
import { IGallery, IGalleryItem } from "../model/gallery.ts";
import { IItem } from "../model/item.ts";

export const GalleryCrudy = new Crudy<IGallery>(`${config.SERVER_URL}/gallery`);

export const GalleryItemCrudy = new Crudy<IGalleryItem>(
  `${config.SERVER_URL}/gallery-item`,
);

export function addItemToGalleries(
  item: IItem["id"],
  galleries: IGallery["id"][],
): Promise<number> {
  return get(
    `${config.SERVER_URL}/gallery-item/save/itemId?itemId=${item}&galleryId=${encodeURIComponent(galleries.join(","))}`,
    {
      method: "POST",
    },
  );
}

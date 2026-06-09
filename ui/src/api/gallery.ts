import Crudy, { antdget, config } from "@allape/gocrud-react";
import { IGallery, IGalleryItem } from "../model/gallery.ts";
import { IItem } from "../model/item.ts";

export const GalleryCrudy = new Crudy<IGallery>(`${config.SERVER_URL}/gallery`);

export function getGalleryItemsByItemIds(
  itemIds: IItem["id"][],
): Promise<IGalleryItem[]> {
  return antdget<IGalleryItem[]>(
    `${config.SERVER_URL}/gallery-item/all?in_itemId=${itemIds.join(",")}`,
  );
}

export function addItemToGalleries(
  itemId: IItem["id"],
  galleryIds: IGallery["id"][],
): Promise<number> {
  return antdget(`${config.SERVER_URL}/gallery-item/save/itemId/${itemId}`, {
    method: "POST",
    body: JSON.stringify(
      galleryIds.map<Pick<IGalleryItem, "itemId" | "galleryId">>((id) => ({
        itemId,
        galleryId: id,
      })),
    ),
  });
}

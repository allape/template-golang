import Crudy, { antdget, config } from "@allape/gocrud-react";
import { IItem, IItemTag } from "../model/item.ts";
import { ITag } from "../model/tag.ts";

export const ItemCrudy = new Crudy<IItem>(`${config.SERVER_URL}/item`);

export function getItemTagsByItemIds(
  itemIds: IItem["id"][],
): Promise<IItemTag[]> {
  return antdget<IItemTag[]>(
    `${config.SERVER_URL}/item-tag/all?in_itemId=${itemIds.join(",")}`,
  );
}

export function addTagsToItem(
  tagIds: ITag["id"][],
  itemId: IItem["id"],
): Promise<IItemTag[]> {
  return antdget(`${config.SERVER_URL}/item-tag/save/itemId/${itemId}`, {
    method: "POST",
    body: JSON.stringify(
      tagIds.map<Pick<IItemTag, "itemId" | "tagId">>((id) => ({
        itemId,
        tagId: id,
      })),
    ),
  });
}

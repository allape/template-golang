import { IBase } from "@allape/gocrud";
import { IItem } from "./item.ts";

export interface IGallery extends IBase {
  name: string;
  isPublic: boolean;
  priority: number;
  description: string;
  createdBy: string;
}

export interface IGalleryItem extends Pick<IBase, "createdAt"> {
  itemId: IItem["id"];
  galleryId: IGallery["id"];
}

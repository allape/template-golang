import { IBase, IBaseSearchParams } from "@allape/gocrud";
import { ITimeSortSearchParams } from "@allape/gocrud/src/model.ts";
import { IItem } from "./item.ts";

export interface IGallery extends IBase {
  name: string;
  isPublic: boolean;
  description: string;
  createdBy: string;
}

export interface IGallerySearchParams
  extends IBaseSearchParams, ITimeSortSearchParams {
  like_name?: string;
  isPublic?: boolean;
  createdBy?: string;
}

export interface IGalleryItem extends Pick<IBase, "createdAt"> {
  itemId: IItem["id"];
  galleryId: IGallery["id"];
}

export interface IGalleryItemSearchParams {
  in_galleryId?: IGalleryItem["galleryId"][];
  in_itemId?: IGalleryItem["itemId"][];
}

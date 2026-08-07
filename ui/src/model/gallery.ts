import { IBase, IBaseSearchParams } from "@allape/gocrud";
import { IItem } from "./item.ts";
import { ITag } from "./tag.ts";

export interface IGallery extends IBase {
  name: string;
  isPublic: boolean;
  description: string;
  createdBy: string;
}

export interface IGallerySearchParams extends IBaseSearchParams {
  like_name?: string;
  isPublic?: boolean;
  createdBy?: string;
}

export interface IGalleryItem extends Pick<IBase, "createdAt"> {
  galleryId: IGallery["id"];
  itemId: IItem["id"];
}

export interface IGalleryTag extends Pick<IBase, "createdAt"> {
  galleryId: IGallery["id"];
  tagId: ITag["id"];
}

import { IBase, IBaseSearchParams } from "@allape/gocrud";
import { IGallery } from "./gallery.ts";
import { ITag } from "./tag.ts";

export interface IItem extends IBase {
  name: string;
  src: string;
  description: string;
  createdBy: string;
}

export interface IItemSearchParams extends IBaseSearchParams {
  in_galleryId?: IGallery["id"][];
  like_name?: string;
  createdBy?: string;
}

export interface IItemTag extends Pick<IBase, "createdAt"> {
  itemId: IItem["id"];
  tagId: ITag["id"];
}

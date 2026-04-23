import { IBase, IBaseSearchParams, SortType } from "@allape/gocrud";
import { ITimeSortSearchParams } from "@allape/gocrud/src/model.ts";
import { IGallery } from "./gallery.ts";
import { ITag } from "./tag.ts";

export interface IItem extends IBase {
  name: string;
  src: string;
  priority: number;
  description: string;
  createdBy: string;
}

export interface IItemSearchParams
  extends IBaseSearchParams,
    Pick<ITimeSortSearchParams, "orderBy_updatedAt"> {
  in_id?: IItem["id"][];
  in_galleryId?: IGallery["id"][];
  like_name?: string;
  createdBy?: string;
  orderBy_priority?: SortType;
  orderByDefault?: string;
}

export interface IItemTag extends Pick<IBase, "createdAt"> {
  itemId: IItem["id"];
  tagId: ITag["id"];
}

export interface IItemTagSearchParams {
  in_itemId?: IItemTag["itemId"][];
  in_tagId?: IItemTag["tagId"][];
}

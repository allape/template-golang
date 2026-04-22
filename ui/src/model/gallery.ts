import { IBase, IBaseSearchParams, SortType } from "@allape/gocrud";
import { ITimeSortSearchParams } from "@allape/gocrud/src/model.ts";

export interface IGallery extends IBase {
  name: string;
  isPublic: boolean;
  priority: number;
  description: string;
  createdBy: string;
}

export interface IGallerySearchParams
  extends IBaseSearchParams, Pick<ITimeSortSearchParams, "orderBy_updatedAt"> {
  in_id?: IGallery["id"][];
  like_name?: string;
  isPublic?: "true" | "false";
  createdBy?: string;
  orderBy_priority?: SortType;
  orderByDefault?: string;
}

import { IBase, IBaseSearchParams, SortType } from "@allape/gocrud";
import { ITimeSortSearchParams } from "@allape/gocrud/src/model.ts";

export interface ITag extends IBase {
  name: string;
  alias: string;
  priority: number;
  description: string;
}

export interface ITagSearchParams
  extends
    IBaseSearchParams,
    Pick<ITimeSortSearchParams, "orderBy_createdAt" | "orderBy_updatedAt"> {
  in_id?: ITag["id"][];
  like_name?: string;
  like_alias?: string;
  orderBy_priority?: SortType;
}

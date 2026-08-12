import { IBase, IBaseSearchParams } from "@allape/gocrud";
import { ITimeSortSearchParams } from "@allape/gocrud/src/model.ts";

export interface ITag extends IBase {
  name: string;
  alias: string;
  color: string;
  description: string;
}

export interface ITagSearchParams
  extends IBaseSearchParams, ITimeSortSearchParams {
  like_name?: string;
  like_alias?: string;
  like_keyword?: string;
}

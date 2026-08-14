import { IBaseSearchParams } from "@allape/gocrud";

export interface ITagSearchParams extends IBaseSearchParams {
  like_name?: string;
  like_alias?: string;
  like_keyword?: string;
}

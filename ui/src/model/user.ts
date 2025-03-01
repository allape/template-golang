import { IBase, IBaseSearchParams } from "@allape/gocrud";

export interface IUser extends IBase {
  name: string;
  description: string;
}

export interface IUserSearchParams extends IBaseSearchParams {
  like_name?: string;
  in_id?: IUser["id"][];
}

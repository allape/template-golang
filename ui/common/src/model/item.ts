import { IBase } from "@allape/gocrud";
import { ITag } from "./tag.ts";

export interface IItem extends IBase {
  name: string;
  src: string;
  thumbnail: string;
  description: string;
  enabled: boolean;
}

export interface IItemTag extends Pick<IBase, "createdAt"> {
  itemId: IItem["id"];
  tagId: ITag["id"];
}

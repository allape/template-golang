import { IBase } from "@allape/gocrud";
import { ITag } from "./tag.ts";

export interface IItem extends IBase {
  name: string;
  src: string;
  priority: number;
  description: string;
  createdBy: string;
}

export interface IItemTag extends Pick<IBase, "createdAt"> {
  itemId: IItem["id"];
  tagId: ITag["id"];
}

import { IBase } from "@allape/gocrud";

export interface ITag extends IBase {
  name: string;
  alias: string;
  color: string;
  description: string;
}

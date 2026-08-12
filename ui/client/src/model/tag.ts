import { IBase } from "@allape/gocrud";

export interface ITag extends IBase {
  name: string;
  alias: string;
  priority: number;
  color: string;
  description: string;
}

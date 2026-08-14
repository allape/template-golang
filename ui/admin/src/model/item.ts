import { IBaseSearchParams } from "@allape/gocrud";
import { IGallery } from "common/src/model/gallery.ts";

export interface IItemSearchParams extends IBaseSearchParams {
  in_galleryId?: IGallery["id"][];
  like_name?: string;
  createdBy?: string;
}

import { IBaseSearchParams } from "@allape/gocrud";
import { IUserGallery } from "common/src/model/user.ts";

/**
 All fields from this type are not present in real data,
 just for the ignorance of type errors
 */
type IFakeIBaseSearchParams = IBaseSearchParams;

export interface IUserGallerySearchParams extends IFakeIBaseSearchParams {
  userId?: IUserGallery["userId"];
  galleryId?: IUserGallery["galleryId"];
}

import { IBase, IBaseSearchParams } from "@allape/gocrud";
import { IGallery } from "./gallery.ts";

/**
 * All fields from this type are not present in real data,
 * just for the ignorance of type errors
 */
type IFakeIBase = IBase;

/**
 * Same as above
 */
type IFakeIBaseSearchParams = IBaseSearchParams;

export interface IUserGallery extends Pick<IBase, "createdAt">, IFakeIBase {
  userId: string;
  galleryId: IGallery["id"];
}

export interface IUserGallerySearchParams extends IFakeIBaseSearchParams {
  userId?: IUserGallery["userId"];
  galleryId?: IUserGallery["galleryId"];
}

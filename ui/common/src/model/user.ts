import { IBase } from "@allape/gocrud";
import { IGallery } from "./gallery.ts";

/**
 * All fields from this type are not present in real data,
 * just for the ignorance of type errors
 */
type IFakeIBase = IBase;

export interface IUserGallery extends Pick<IBase, "createdAt">, IFakeIBase {
  userId: string;
  galleryId: IGallery["id"];
}

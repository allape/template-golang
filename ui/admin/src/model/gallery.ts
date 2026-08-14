import { IBaseSearchParams } from "@allape/gocrud";

export interface IGallerySearchParams extends IBaseSearchParams {
  like_name?: string;
  isPublic?: boolean;
  createdBy?: string;
}

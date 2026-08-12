import AntdCrudy, { antdget, config } from "@allape/gocrud-react";
import { IUserGallery, IUserGallerySearchParams } from "../model/user.ts";

export const UserGalleryCrudy = new AntdCrudy<
  IUserGallery,
  IUserGallerySearchParams
>(`${config.SERVER_URL}/user-gallery`);

export function saveAllUserGallery(
  userGalleries: Array<Partial<IUserGallery> & Pick<IUserGallery, "userId" | "galleryId">>
): Promise<number> {
  return antdget(`${config.SERVER_URL}/user-gallery/all`, {
    method: "PUT",
    body: JSON.stringify(userGalleries)
  });
}

export function deleteUserGallery(
  userId: IUserGallery["userId"],
  galleryId: IUserGallery["galleryId"]
): Promise<void> {
  return antdget(
    `${config.SERVER_URL}/user-gallery/0?userId=${userId}&galleryId=${galleryId}`,
    {
      method: "DELETE"
    }
  );
}

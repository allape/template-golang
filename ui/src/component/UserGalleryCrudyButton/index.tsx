import {
  asDefaultPattern,
  CrudyButton,
  ICrudyButtonProps,
} from "@allape/gocrud-react";
import { Form, InputNumber, TableColumnsType } from "antd";
import { ReactElement, useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { GalleryCrudy } from "../../api/gallery.ts";
import {
  deleteUserGallery,
  saveAllUserGallery,
  UserGalleryCrudy,
} from "../../api/user.ts";
import { IGallery } from "../../model/gallery.ts";
import { IUserGallery, IUserGallerySearchParams } from "../../model/user.ts";
import GallerySelector from "../GallerySelector";

interface IUserGalleryModified extends IUserGallery {
  _fakeId?: string;

  _gallery?: IGallery;

  _galleryIds?: IUserGallery["galleryId"][];
}

type IRecord = IUserGalleryModified;
type ISearchParams = IUserGallerySearchParams;

const crudy = UserGalleryCrudy;

export type IUserGalleryCrudyButtonProps = Partial<ICrudyButtonProps<IRecord>>;

export default function UserGalleryCrudyButton(
  props: IUserGalleryCrudyButtonProps,
): ReactElement {
  const { t } = useTranslation();

  const [searchParams] = useState<ISearchParams>(() => ({}));

  const columns = useMemo<TableColumnsType<IRecord>>(
    () => [
      {
        title: t("user.gallery.userId"),
        dataIndex: "userId",
      },
      {
        title: t("user.gallery.gallery"),
        dataIndex: "galleryId",
        render: (v: IRecord["galleryId"], record) =>
          `${v}: ${record._gallery?.name}`,
      },
      {
        title: t("createdAt"),
        dataIndex: "createdAt",
        render: asDefaultPattern,
      },
    ],
    [t],
  );

  const handleAfterListed = useCallback(
    async (records: IRecord[]): Promise<IRecord[]> => {
      if (records.length === 0) {
        return [];
      }

      records.forEach((record) => {
        record._fakeId = `${record.userId}:${record.galleryId}`;
      });

      const galleryIds = records.map((record) => record.galleryId);
      if (galleryIds.length > 0) {
        const galleries = await GalleryCrudy.all({
          in_id: galleryIds,
        });

        records.forEach((record) => {
          const gallery = galleries.find((g) => g.id === record.galleryId);
          if (!gallery) {
            return;
          }
          record._gallery = gallery;
        });
      }

      return records;
    },
    [],
  );

  const handleSave = useCallback(async (record: IRecord): Promise<IRecord> => {
    if (record._galleryIds?.length) {
      await saveAllUserGallery(
        record._galleryIds.map((gi) => ({
          userId: `${record.userId}`,
          galleryId: gi,
        })),
      );
    }
    return record;
  }, []);

  const handleDelete = useCallback(async (record: IRecord) => {
    return deleteUserGallery(record.userId, record.galleryId);
  }, []);

  return (
    <CrudyButton<IRecord, ISearchParams>
      name={t("user.gallery._")}
      columns={columns}
      crudy={crudy}
      searchParams={searchParams}
      afterListed={handleAfterListed}
      editable={false}
      onSave={handleSave}
      onDelete={handleDelete}
      {...props}
    >
      <Form.Item name="userId" label={t("user.gallery.userId")}>
        <InputNumber
          precision={0}
          step={1}
          min={1}
          max={Number.MAX_SAFE_INTEGER}
          placeholder={t("user.gallery.userId")}
        />
      </Form.Item>
      <Form.Item name="_galleryIds" label={t("user.gallery.gallery")}>
        <GallerySelector mode="multiple" />
      </Form.Item>
    </CrudyButton>
  );
}

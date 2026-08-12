import { useLoading, useProxy } from "@allape/use-loading";
import { Button } from "antd";
import { ReactElement, useCallback, useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router";
import {
  getDetailById,
  IGalleryDetail,
  toImageURL,
} from "../../api/gallery.ts";
import { IItem } from "../../model/item.ts";
import { ITag } from "../../model/tag.ts";

interface IItemModified extends IItem {
  _src?: string;
  _tags?: ITag[];
}

interface IGalleryDetailModified extends IGalleryDetail {
  _items?: IItemModified[];
}

export default function Gallery(): ReactElement {
  const { t } = useTranslation();
  const { loading, execute } = useLoading();

  const { galleryId } = useParams<Record<"galleryId", string | undefined>>();

  const stoperRef = useRef<AbortController>();

  const [detail, , setDetail] = useProxy<IGalleryDetailModified | undefined>(undefined);
  const [items, setItems] = useState<IItemModified[]>([]);

  const getDetail = useCallback(() => {
    if (!galleryId) {
      return;
    }

    const parsedId = parseInt(galleryId);
    if (Number.isNaN(parsedId)) {
      return;
    }

    if (stoperRef.current) {
      stoperRef.current.abort();
    }
    stoperRef.current = new AbortController();

    execute(async () => {
      const detail = await getDetailById(parsedId, {
        signal: stoperRef.current?.signal,
      });
      setDetail(detail);
      setItems(
        detail.items.map((item) => ({
          ...item,
          _src: toImageURL(parsedId, item.id),
          _tags: detail.itemTags
            .filter((it) => it.itemId === item.id)
            .map((it) => detail.tags.find((t) => t.id === it.tagId) as ITag)
            .filter(Boolean),
        })),
      );
    }).then();
  }, [execute, galleryId, setDetail]);

  useEffect(() => {
    const timerId = setTimeout(() => {
      getDetail();
    }, 100);
    return () => {
      clearTimeout(timerId);
    };
  }, [execute, getDetail]);

  return (
    <div>
      {loading && <Button loading type="link" />}
      {!galleryId ? <div>{t("gallery.galleryNotValid")}</div> : undefined}
      {detail && <div>{detail.gallery.name}</div>}
      {items.map((item) => (
        <div key={item.id}>
          <img style={{ width: "100%" }} src={item._src} alt={item.name} />
        </div>
      ))}
    </div>
  );
}

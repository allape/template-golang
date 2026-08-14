import { useLoading, useProxy } from "@allape/use-loading";
import { HomeOutlined, LoadingOutlined } from "@ant-design/icons";
import { Button, Tag } from "antd";
import cls from "classnames";
import { IItem } from "common/src/model/item.ts";
import { ITag } from "common/src/model/tag.ts";
import { ReactElement, useCallback, useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router";
import {
  getDetailById,
  IGalleryDetail,
  toImageURL,
} from "../../api/gallery.ts";
import styles from "./style.module.scss";

interface IItemModified extends IItem {
  _src?: string;
  _tags?: ITag[];
}

interface IGalleryDetailModified extends IGalleryDetail {
  _items?: IItemModified[];
}

export default function Gallery(): ReactElement {
  const navigation = useNavigate();
  const { t } = useTranslation();
  const { execute } = useLoading();

  const { galleryId } = useParams<Record<"galleryId", string | undefined>>();

  const stoperRef = useRef<AbortController>();

  const [detail, , setDetail] = useProxy<IGalleryDetailModified | undefined>(
    undefined,
  );
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

  const handleBackHome = useCallback(() => {
    navigation("/");
  }, [navigation]);

  const openImageInNewTag = useCallback((item?: IItemModified) => {
    if (!item?._src) {
      return;
    }
    window.open(item._src, "_blank");
  }, []);

  return (
    <div className={styles.wrapper}>
      <div className={styles.title}>
        <Button
          className={styles.homeButton}
          type="link"
          size="large"
          onClick={handleBackHome}
        >
          <HomeOutlined />
        </Button>
        <div className={styles.text}>
          <div className={styles.name}>
            {detail?.gallery.name || <LoadingOutlined />}
          </div>
          <div className={styles.tags}>
            {detail?.tags?.map((tag) => (
              <Tag key={tag.id} className={styles.tag}>
                {tag.name}
              </Tag>
            ))}
          </div>
        </div>
        <Button
          className={cls(styles.homeButton, styles.hidden)}
          type="link"
          size="large"
          onClick={handleBackHome}
        >
          <HomeOutlined />
        </Button>
      </div>

      {!galleryId ? <div>{t("gallery.galleryNotValid")}</div> : undefined}

      <div className={styles.items}>
        {items.map((item) => (
          <div
            key={item.id}
            className={styles.item}
            onClick={() => openImageInNewTag(item)}
          >
            <img
              className={styles.cover}
              loading="lazy"
              src={item._src}
              alt={item.name}
            />
            <div className={styles.tags}>
              {item._tags?.map((tag) => (
                <Tag key={tag.id} className={styles.tag}>
                  {tag.name}
                </Tag>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

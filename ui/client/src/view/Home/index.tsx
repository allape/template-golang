import { useLoading, useProxy } from "@allape/use-loading";
import { LoadingOutlined } from "@ant-design/icons";
import { Empty } from "antd";
import { ReactElement, useEffect } from "react";
import { useNavigate } from "react-router";
import {
  getAllGalleries,
  IGalleryInfo,
  toImageURL,
} from "../../api/gallery.ts";
import styles from "./style.module.scss";

interface IGalleryInfoModified extends IGalleryInfo {
  _coverLink?: string;
}

export default function Home(): ReactElement {
  const navigate = useNavigate();
  const { loading, execute } = useLoading();

  const [galleries, , setGalleries] = useProxy<IGalleryInfoModified[]>([]);

  useEffect(() => {
    execute(async () => {
      const infos: IGalleryInfoModified[] = await getAllGalleries();
      infos.forEach((info) => {
        info._coverLink = toImageURL(info.info.id);
      });
      setGalleries(infos);
    }).then();
  }, [execute, setGalleries]);

  return (
    <div className={styles.wrapper}>
      {loading ? (
        <div className={styles.loading}>
          <LoadingOutlined />
        </div>
      ) : undefined}
      {galleries.length === 0 ? (
        <Empty />
      ) : (
        galleries.map((gallery) => (
          <div
            key={gallery.info.id}
            className={styles.gallery}
            title={gallery.info.name}
            onClick={() => navigate(`/gallery/${gallery.info.id}`)}
          >
            <img
              className={styles.cover}
              loading="lazy"
              src={gallery._coverLink}
              alt={gallery.info.name}
            />
            <div className={styles.name}>{gallery.info.name}</div>
          </div>
        ))
      )}
    </div>
  );
}

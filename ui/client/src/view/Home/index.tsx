import { useLoading, useProxy } from "@allape/use-loading";
import { Empty, Spin } from "antd";
import { ReactElement, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import {
  getAllGalleries,
  IGalleryInfo,
  toImageURL,
} from "../../api/gallery.ts";

interface IGalleryInfoModified extends IGalleryInfo {
  _coverLink?: string;
}

export default function Home(): ReactElement {
  const navigate = useNavigate();
  const { t } = useTranslation();
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
    <div>
      <Spin spinning={loading}>
        {galleries.length === 0 ? <Empty /> : undefined}
        {galleries.map((gallery) => (
          <div
            key={gallery.info.id}
            onClick={() => navigate(`/gallery/${gallery.info.id}`)}
          >
            <img width={100} src={gallery._coverLink} alt={gallery.info.name} />
            <div title={t("gallery.name")}>{gallery.info.name}</div>
          </div>
        ))}
      </Spin>
    </div>
  );
}

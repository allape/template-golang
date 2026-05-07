import { useLoading, useProxy } from "@allape/use-loading";
import { Empty, Spin } from "antd";
import { ReactElement, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { getAllGalleries, toImageURL } from "../../api/gallery.ts";
import { IGallery } from "../../model/gallery.ts";

export default function Home(): ReactElement {
  const { t } = useTranslation();
  const { loading, execute } = useLoading();

  const [galleries, , setGalleries] = useProxy<IGallery[]>([]);

  useEffect(() => {
    execute(async () => {
      const gal = await getAllGalleries();
      setGalleries(gal);
    }).then();
  }, [execute, setGalleries]);

  // a lot things todo

  return (
    <div>
      <Spin spinning={loading}>
        {galleries.length === 0 ? <Empty /> : undefined}
        {galleries.map((gallery) => (
          <div key={gallery.id}>
            <img width={100} src={toImageURL(gallery.id)} alt={gallery.name} />
            <div title={t("gallery.name")}>{gallery.name}</div>
          </div>
        ))}
      </Spin>
    </div>
  );
}

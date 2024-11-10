import Album from "../../types/models/Album.ts";
import {Image} from "@mantine/core";

interface AlbumCardProps  {
  album: Album
}

function AlbumCard({ album } : AlbumCardProps) {
  return (
    <div>
      <Image src={album.imageUrl} />
      {album.title}
    </div>
  );
}

export default AlbumCard;

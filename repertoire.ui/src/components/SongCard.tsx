import React from 'react';
import Song from "../types/models/Song";

interface SongCardProps {
  song: Song
}

function SongCard({ song }: SongCardProps) {
  return (
    <div>
      <h2>{song.title}</h2>
      <div>{song.isRecorded}</div>
    </div>
  );
}

export default SongCard;

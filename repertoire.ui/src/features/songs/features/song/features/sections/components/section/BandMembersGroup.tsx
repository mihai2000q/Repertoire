import { Avatar, AvatarGroup } from '@mantine/core'
import { IconUser } from '@tabler/icons-react'
import { BandMember } from '../../../../../../../../types/models/Artist.ts'

interface BandMembersGroupProps {
  bandMembers: BandMember[]
}

function BandMembersGroup({ bandMembers }: BandMembersGroupProps) {
  return (
    <AvatarGroup>
      {bandMembers.map((bandMember) => (
        <Avatar
          key={bandMember.id}
          size={'28px'}
          bd={0}
          color={bandMember.color}
          src={bandMember.imageUrl}
          alt={bandMember.imageUrl && bandMember.name}
        >
          <IconUser aria-label={`default-icon-${bandMember.name}`} size={15} />
        </Avatar>
      ))}
    </AvatarGroup>
  )
}

export default BandMembersGroup

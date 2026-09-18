import { Avatar, AvatarGroup, AvatarGroupProps } from '@mantine/core'
import { IconUser } from '@tabler/icons-react'
import { BandMember } from '../../types/models/Artist.ts'

interface BandMembersGroupProps extends AvatarGroupProps {
  bandMembers: BandMember[]
}

function BandMembersGroup({ bandMembers, ...props }: BandMembersGroupProps) {
  return (
    <AvatarGroup {...props}>
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

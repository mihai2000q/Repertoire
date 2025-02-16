import {http, HttpResponse} from "msw";
import {setupServer} from "msw/node";
import {reduxRender} from "../../../test-utils.tsx";
import BandMembersCard from "./BandMembersCard.tsx";
import {BandMember} from "../../../types/models/Artist.ts";
import {screen} from "@testing-library/react";
import {userEvent} from "@testing-library/user-event";

describe('Band Members Card', () => {
  const handlers = [
    http.get('/artists/band-members/roles', async () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const bandMembers: BandMember[] = [
    {
      id: '1',
      name: 'Member 1',
      color: '#123123',
      roles: []
    },
    {
      id: '2',
      name: 'Member 2',
      roles: []
    }
  ]

  it('should render with band members', async () => {
    const user = userEvent.setup()

    reduxRender(<BandMembersCard bandMembers={bandMembers} artistId={''} />)

    expect(screen.getByText(/band members/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'back' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'forward' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'band-members-more-menu' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'band-members-more-menu' }))
    expect(screen.getByRole('menuitem', { name: /add new member/i })).toBeInTheDocument()

    for (const bandMember of bandMembers) {
      expect(screen.getByLabelText(`band-member-card-${bandMember.name}`))
    }
    expect(screen.queryByLabelText('add-new-band-member-card')).not.toBeInTheDocument()
  })

  it('should render accordingly without members', async () => {
    const user = userEvent.setup()

    reduxRender(<BandMembersCard bandMembers={[]} artistId={''} />)

    expect(screen.getByText(/band members/i)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'back' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'forward' })).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'band-members-more-menu' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'band-members-more-menu' }))
    expect(screen.getByRole('menuitem', { name: /add new member/i })).toBeInTheDocument()

    expect(screen.getByLabelText('add-new-band-member-card')).toBeInTheDocument()
  })

  it.skip('should be able to reorder', () => {
    reduxRender(<BandMembersCard bandMembers={[]} artistId={''} />)
  })

  describe('on menu', () => {
    it('should open add new member modal when clicking on add new member', async () => {
      const user = userEvent.setup()

      reduxRender(<BandMembersCard bandMembers={[]} artistId={''} />)

      await user.click(screen.getByRole('button', { name: 'band-members-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add new member/i }))
      expect(await screen.findByRole('dialog', { name: /add new band member/i })).toBeInTheDocument()
    })
  })
})

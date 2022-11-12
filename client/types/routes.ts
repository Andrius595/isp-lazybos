const Routes = {
  Main: 'main',
  Profile: 'profile',
  Auth: {
    Login: 'login',
    Register: 'register',
  },
  Identity: {
    List: 'identity',
    Request: 'identity-request',
    Confirm: 'identity-id',
  },
  Matches: {
    List: 'matches',
    BetForm: 'matches-bet-form',
    Information: 'matches-id'
  },
  AdminMatches: {
    List: 'admin-matches',
    RemoveConfirmation: 'admin-matches-remove-confirmation',
    EndConfirmation: 'admin-matches-end-confirmation',
    Create: 'admin-matches-create'
  },
  Bets: {
    List: 'bets',
  },
  AutomaticBets: {
    List: 'automatic-bets',
    Create: 'automatic-bets-create',
  },
  Wallet: {
    Deposit: 'wallet-deposit',
    Activate: 'wallet-activate',
    Withdraw: 'wallet-withdraw',
  },
  Reports: {
    List: 'reports',
    Profit: 'reports-profit',
    Taxes: 'reports-taxes',
    Matches: 'reports-matches',
    Users: 'reports-users',
    AdminsActions: 'reports-admins-actions',
    ScheduledReport: 'reports-scheduled-report',
  }
}

export default Routes
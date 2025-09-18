export function jumpPlayerList(user_id) {
    const { href } = this.$router.resolve({ name: 'player_info_list', query: { user_id: user_id } })
    window.open(href, '_blank');
  }

export function jumpSummaryRevenueConfig(date, channel) {
  const { href } = this.$router.resolve({ name: 'summary_revenue_config', query: {date: date, channel: channel} })
  window.open(href, '_blank');
}

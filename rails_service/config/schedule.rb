# Use this file to easily define all of your cron jobs.
#
# It's helpful, but not entirely necessary to understand cron before proceeding.
# http://en.wikipedia.org/wiki/Cron
ENV.each_key do |key|
  env key.to_sym, ENV[key]
end
set :output, "log/cron_log.log"
set :environment, ENV['RAILS_ENV']


every 1.minute do
  runner 'PersistCountersJob.perform_now'
end

# Example:
#
#
# every 2.hours do
#   command "/usr/bin/some_great_command"
#   runner "MyModel.some_method"
#   rake "some:great:rake:task"
# end
#
# every 4.days do
#   runner "AnotherModel.prune_old_records"
# end

# Learn more: http://github.com/javan/whenever

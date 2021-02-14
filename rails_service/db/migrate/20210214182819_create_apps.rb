class CreateApps < ActiveRecord::Migration[5.0]
  def change
    create_table :apps do |t|
      t.string :token
      t.string :name
      t.integer :chats_count, default: 0
      t.timestamps
    end
    add_index :apps, :token
  end
end

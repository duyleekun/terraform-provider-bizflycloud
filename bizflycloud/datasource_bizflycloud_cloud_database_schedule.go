// This file is part of terraform-provider-bizflycloud
//
// Copyright (C) 2021  Bizfly Cloud
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>

package bizflycloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizFlyCloudDatabaseSchedule() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizFlyCloudDatabaseScheduleRead,
		Schema: resourceCloudDatabaseScheduleSchema(),
	}
}

func dataSourceBizFlyCloudDatabaseScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	scheduleID := d.Id()

	log.Printf("[DEBUG] Reading Database schedule: %s", scheduleID)
	schedule, err := client.CloudDatabase.Schedules().Get(context.Background(), scheduleID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing Database schedule: %w", err)
	}

	log.Printf("[DEBUG] Found Database schedule: %s", scheduleID)
	log.Printf("[DEBUG] bizflycloud_cloud_database_schedule - Single Database schedule found: %s", schedule.Name)

	d.SetId(schedule.ID)

	return nil
}

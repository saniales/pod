Vue.component('PageSettings', {
	name: 'Settings',
		data () { return { 
		duOSsettings }},
		template: `<main class="pageSettings"><div class="rwrap">
			<div v-html="duOSsettings.data.daemon.schema"></div>
			<vue-form-generator class="flx flc fii" :schema="duOSsettings.data.daemon.schema" :model="duOSsettings.data.daemon.config"></vue-form-generator>
		</div></main>`
	});
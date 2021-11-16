<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.5" tiledversion="1.7.2" name="tileset" tilewidth="16" tileheight="16" tilecount="64" columns="8">
 <image source="tileset.png" width="128" height="128"/>
 <tile id="0">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="16" height="16"/>
  </objectgroup>
 </tile>
 <tile id="1">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="16" height="16"/>
  </objectgroup>
 </tile>
 <tile id="2">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="16" height="16"/>
  </objectgroup>
 </tile>
 <tile id="3">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="15" height="15"/>
  </objectgroup>
 </tile>
 <tile id="8">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="16" height="16"/>
  </objectgroup>
 </tile>
 <tile id="10">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="16" height="16"/>
  </objectgroup>
 </tile>
 <tile id="11">
  <objectgroup draworder="index" id="2">
   <object id="1" x="7" y="8" width="9" height="8"/>
  </objectgroup>
 </tile>
 <tile id="12">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="8" width="9" height="8"/>
  </objectgroup>
 </tile>
 <tile id="13">
  <animation>
   <frame tileid="13" duration="150"/>
   <frame tileid="14" duration="150"/>
   <frame tileid="15" duration="150"/>
  </animation>
 </tile>
 <tile id="14">
  <animation>
   <frame tileid="15" duration="150"/>
   <frame tileid="13" duration="150"/>
   <frame tileid="14" duration="150"/>
  </animation>
 </tile>
 <tile id="15">
  <animation>
   <frame tileid="15" duration="150"/>
   <frame tileid="14" duration="150"/>
   <frame tileid="13" duration="150"/>
  </animation>
 </tile>
 <tile id="16">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="16" height="16"/>
  </objectgroup>
 </tile>
 <tile id="17">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="16" height="16"/>
  </objectgroup>
 </tile>
 <tile id="18">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="16" height="16"/>
  </objectgroup>
 </tile>
 <tile id="19">
  <objectgroup draworder="index" id="2">
   <object id="1" x="7" y="0" width="9" height="8"/>
  </objectgroup>
 </tile>
 <tile id="20">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="9" height="8"/>
  </objectgroup>
 </tile>
 <wangsets>
  <wangset name="Cave" type="corner" tile="-1">
   <wangcolor name="Wall" color="#ff0000" tile="-1" probability="1"/>
   <wangtile tileid="0" wangid="0,0,0,1,0,0,0,0"/>
   <wangtile tileid="1" wangid="0,0,0,1,0,1,0,0"/>
   <wangtile tileid="2" wangid="0,0,0,0,0,1,0,0"/>
   <wangtile tileid="8" wangid="0,1,0,1,0,0,0,0"/>
   <wangtile tileid="10" wangid="0,0,0,0,0,1,0,1"/>
   <wangtile tileid="11" wangid="0,1,0,0,0,1,0,1"/>
   <wangtile tileid="12" wangid="0,1,0,1,0,0,0,1"/>
   <wangtile tileid="16" wangid="0,1,0,0,0,0,0,0"/>
   <wangtile tileid="17" wangid="0,1,0,0,0,0,0,1"/>
   <wangtile tileid="18" wangid="0,0,0,0,0,0,0,1"/>
   <wangtile tileid="19" wangid="0,0,0,1,0,1,0,1"/>
   <wangtile tileid="20" wangid="0,1,0,1,0,1,0,0"/>
   <wangtile tileid="39" wangid="0,1,0,1,0,1,0,1"/>
  </wangset>
 </wangsets>
</tileset>

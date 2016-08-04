## finder_machine

Automatiza la transferencia de comprobantes entre diferentes orígenes de datos, como bases de datos, colas y sistema de archivos.

## Setup

    $ go get github.com/edmt/finder_machine

## Stored Procedures

FinderMachine_ReadXml:

```
CREATE PROCEDURE FinderMachine_ReadXml
    @startDate VARCHAR(20),
    @endDate varchar(20)
AS
BEGIN

SET NOCOUNT ON

select xml.uuid, xml.xml, xml.timestamp
from xml
inner join [QA-DANY].dbo.cfd on cfd.numTimbre = xml.uuid
left join  [QA-DANY].dbo.POOL_ENVIOCFD_SAT_Z as pz on pz.comprobante_Id = cfd.idInternal
left join  [QA-DANY].dbo.cfd_envio_sat_z as acuse on acuse.comprobante_Id = cfd.idInternal
where timestamp > @startDate and timestamp < @endDate and pz.comprobante_Id is null and acuse.comprobante_Id is null

END


exec FinderMachine_ReadXML '2016-05-16', '2016-05-17'
```

FinderMachine_WritePool:

```
CREATE PROCEDURE FinderMachine_WritePool
    @uuid VARCHAR(40)
AS
BEGIN

SET NOCOUNT ON

insert into [QA-DANY].dbo.POOL_ENVIOCFD_SAT_Z(idInternal, comprobante_Id, fechaRegistro, status)
select      replace(newid(), '-', ''), idInternal, getdate(), 0
from        [QA-DANY].dbo.cfd
where numtimbre = @uuid;

END

exec FinderMachine_WritePool 'fc0d9501-25f9-40fa-b4ba-73dfaf06dc6d'
```

FinderMachine_ReadXml_MissingCfd:

```
CREATE PROCEDURE FinderMachine_ReadXml_MissingCfd
    @startDate VARCHAR(20),
    @endDate varchar(20)
AS
BEGIN

SET NOCOUNT ON

select xml.uuid, xml.xml, xml.timestamp
from xml
left join [QA-DANY].dbo.cfd  as cfd0 on cfd0.numTimbre = xml.uuid
left join [QA-BF].dbo.cfd    as cfd1 on cfd1.numTimbre = xml.uuid
left join [QA-TF].dbo.cfd    as cfd2 on cfd2.numTimbre = xml.uuid
left join [QA-CF].dbo.cfd    as cfd3 on cfd3.numTimbre = xml.uuid
where timestamp > @startDate and timestamp < @endDate
and cfd0.idInternal is null
and cfd1.idInternal is null
and cfd2.idInternal is null
and cfd3.idInternal is null
order by timestamp desc

END

exec FinderMachine_ReadXml_MissingCfd '2016-01-01', '2016-06-30'
```


FinderMachine_RecoverDeletedCfd:

```
CREATE PROCEDURE FinderMachine_RecoverDeletedCfd
    @uuid VARCHAR(32)
AS
BEGIN

SET NOCOUNT ON

IF OBJECT_ID(  '[QA-DANY].dbo.cfd_delete', 'U') IS NOT NULL
begin
  insert        [QA-DANY].dbo.cfd
  select * from [QA-DANY].dbo.cfd_delete where numTimbre = @uuid;
end

IF OBJECT_ID(  '[QA-BF].dbo.cfd_delete', 'U') IS NOT NULL
begin
  insert        [QA-BF].dbo.cfd
  select * from [QA-BF].dbo.cfd_delete where numTimbre = @uuid;
end

IF OBJECT_ID(  '[QA-TF].dbo.cfd_delete', 'U') IS NOT NULL
begin
  insert        [QA-TF].dbo.cfd
  select * from [QA-TF].dbo.cfd_delete where numTimbre = @uuid;
end

IF OBJECT_ID(  '[QA-CF].dbo.cfd_delete', 'U') IS NOT NULL
begin
  insert        [QA-CF].dbo.cfd
  select * from [QA-CF].dbo.cfd_delete where numTimbre = @uuid;
end

END

exec FinderMachine_RecoverDeletedCfd 'f49f1605-7a37-4223-a2e2-941a9152e0bf'
```
